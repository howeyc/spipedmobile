package com.howeyc.spiped;

import android.app.Activity;
import android.os.Bundle;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.webkit.WebChromeClient;

import android.webkit.ValueCallback;
import android.net.Uri;
import android.content.Intent;

import android.view.Menu;
import android.view.MenuItem;

import go.Go;
import go.spipedmobile.Spipedmobile;


public class MyActivity extends Activity {

    WebView browser;

    private ValueCallback<Uri> mUploadMessage;  
    private final static int FILECHOOSER_RESULTCODE=1;

    @Override  
    protected void onActivityResult(int requestCode, int resultCode, Intent intent) {  
        if(requestCode==FILECHOOSER_RESULTCODE)
        {  
            if (null == mUploadMessage) return;  
            Uri result = intent == null || resultCode != RESULT_OK ? null : intent.getData();  
            mUploadMessage.onReceiveValue(result);  
            mUploadMessage = null;  
        }
    }

    public static final int MENU_PIPES = Menu.FIRST;
    public static final int MENU_LICENSE = Menu.FIRST + 1;

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        super.onCreateOptionsMenu(menu);

        menu.add(Menu.NONE, MENU_PIPES, Menu.NONE, "Current Pipes");
        menu.add(Menu.NONE, MENU_LICENSE, Menu.NONE, "License Info");
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item)
    {
        switch(item.getItemId())
        {
            case MENU_PIPES:
            browser.loadUrl("http://localhost:56056");
            return true;
            case MENU_LICENSE:
            browser.loadUrl("http://localhost:56056/license.html");
            return true;
        default:
            return super.onOptionsItemSelected(item);
        }
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.web_view);

        //Magic
        Go.init(getApplicationContext());

        Spipedmobile.Start();

        browser = (WebView)findViewById(R.id.browser);
        browser.getSettings().setJavaScriptEnabled(true);
        browser.setWebViewClient(new WebViewClient());
        browser.setWebChromeClient(new WebChromeClient()
        {
        //The undocumented magic method override  
        //Eclipse will swear at you if you try to put @Override here  
        // For Android 3.0+
        public void openFileChooser(ValueCallback<Uri> uploadMsg) {  
            mUploadMessage = uploadMsg;  
            Intent i = new Intent(Intent.ACTION_GET_CONTENT);  
            i.addCategory(Intent.CATEGORY_OPENABLE);  
            i.setType("*/*");  
            MyActivity.this.startActivityForResult(Intent.createChooser(i,"File Chooser"), FILECHOOSER_RESULTCODE);  
        }

        // For Android 3.0+
        public void openFileChooser(ValueCallback uploadMsg, String acceptType) {
            mUploadMessage = uploadMsg;
            Intent i = new Intent(Intent.ACTION_GET_CONTENT);
            i.addCategory(Intent.CATEGORY_OPENABLE);
            i.setType("*/*");
            MyActivity.this.startActivityForResult(Intent.createChooser(i, "File Browser"), FILECHOOSER_RESULTCODE);
        }

        //For Android 4.1
        public void openFileChooser(ValueCallback<Uri> uploadMsg, String acceptType, String capture){
            mUploadMessage = uploadMsg;  
            Intent i = new Intent(Intent.ACTION_GET_CONTENT);  
            i.addCategory(Intent.CATEGORY_OPENABLE);  
            i.setType("image/*");  
            MyActivity.this.startActivityForResult( Intent.createChooser( i, "File Chooser" ), MyActivity.FILECHOOSER_RESULTCODE );
        }
        });
        browser.loadUrl("http://localhost:56056");
    }
}
